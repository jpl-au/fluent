package audio

// Byte constants for HTML rendering.
var (
	TagOpen  = []byte("<audio")
	TagClose = []byte("</audio>")

	AttrAutoplay              = []byte(" autoplay")
	AttrControls              = []byte(" controls")
	AttrLoop                  = []byte(" loop")
	AttrMuted                 = []byte(" muted")
	AttrSrc                   = []byte(" src=\"")
	AttrControlsList          = []byte(" controlslist=\"")
	AttrCrossOrigin           = []byte(" crossorigin=\"")
	AttrDisableRemotePlayback = []byte(" disableremoteplayback")
	AttrLoading               = []byte(" loading=\"")
	AttrPreload               = []byte(" preload=\"")
)
