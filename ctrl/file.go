package ctrl

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"github.com/zicare/rgm/ds"
	"github.com/zicare/rgm/msg"
)

type FileController struct {
	DS      ds.IFileDataSource
	MaxSize int64
}

func (fc FileController) SaveBlob(c *gin.Context, scope ds.StoragePath) (*ds.FileInfo, error) {

	if fc.DS == nil {
		c.JSON(http.StatusInternalServerError, msg.Get("25").SetArgs("file store not configured"))
		return nil, &ds.InvalidFileError{Message: msg.Get("25").SetArgs("file store not configured")}
	}

	hdr, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, msg.Get("missing_file"))
		return nil, err
	}

	if fc.MaxSize > 0 && hdr.Size > fc.MaxSize {
		c.JSON(http.StatusRequestEntityTooLarge, msg.Get("file_too_large"))
		return nil, &ds.InvalidFileError{Message: msg.Get("file_too_large")}
	}

	src, err := hdr.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, msg.Get("cannot_open_file"))
		return nil, err
	}
	defer src.Close()

	buf := make([]byte, 512)
	n, err := io.ReadFull(src, buf)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		c.JSON(http.StatusInternalServerError, msg.Get("cannot_read_file"))
		return nil, err
	}

	input := ds.UploadInput{
		OriginalName: hdr.Filename,
		ContentType:  http.DetectContentType(buf[:n]),
		SizeBytes:    hdr.Size,
		Body: io.MultiReader(
			bytes.NewReader(buf[:n]),
			src,
		),
	}

	info, err := fc.DS.Upload(c.Request.Context(), input, scope)
	if err != nil {
		switch err.(type) {

		case *ds.NotFoundFileError:
			c.JSON(http.StatusNotFound, err)

		case *ds.InvalidFileError:
			c.JSON(http.StatusBadRequest, err)

		default:
			c.JSON(
				http.StatusInternalServerError,
				msg.Get("25").SetArgs(err.Error()),
			)
		}
		return nil, err
	}

	return &info, nil
}

func (fc FileController) SaveMeta(c *gin.Context, d ds.IDataSource, info ds.FileInfo) {

	if qo, err := ds.QOFactory(c, d); err != nil {

		if e := fc.DS.Delete(c.Request.Context(), info.Path); e != nil {
			glog.Errorf("file rollback delete failed: path=%s err=%T(%v)", info.Path, e, e)
		}

		c.JSON(
			http.StatusInternalServerError,
			msg.Get("25").SetArgs(fmt.Sprintf("%T", err), err.Error()),
		)

	} else if err := d.Insert(qo); err != nil {

		if e := fc.DS.Delete(c.Request.Context(), info.Path); e != nil {
			glog.Errorf("file rollback delete failed: path=%s err=%T(%v)", info.Path, e, e)
		}

		switch err.(type) {
		case *ds.ForeignKeyConstraint:
			c.JSON(
				http.StatusConflict,
				msg.Get("42"),
			)
		default:
			c.JSON(
				http.StatusInternalServerError,
				msg.Get("25").SetArgs(fmt.Sprintf("%T", err), err.Error()),
			)
		}

	} else {

		c.JSON(
			http.StatusCreated,
			d,
		)
	}
}
