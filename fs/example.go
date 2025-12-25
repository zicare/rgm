package fs

// Documentation-only usage sketch (not referenced by code):
//
// store, _ := fs.New(fs.Opts{
//     BaseDir: "/var/app/files",
//     PathFn:  fs.Shard2PathFn, // optional
// })
//
// ctrl := ctrl.FileController{
//     DS:      store,
//     MaxSize: 25 << 20,
// }
