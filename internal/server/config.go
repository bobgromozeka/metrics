package server

type StartupConfig struct {
	ServerAddr      string
	StoreInterval   uint
	FileStoragePath string
	Restore         bool
	DatabaseDsn     string
}
