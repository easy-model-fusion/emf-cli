package app

const Name = "emf-cli"

var (
	// Version is the binary version + build number
	Version string
	// BuildDate is the date of build
	BuildDate string
)

func Init(version, buildDate string) {
	Version = version
	BuildDate = buildDate
	initLogger()
}
