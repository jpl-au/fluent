package video

// Byte constants for HTML rendering.
var (
	TagOpen  = []byte("<video")
	TagClose = []byte("</video>")

	AttrSrc                     = []byte(" src=\"")
	AttrAutoplay                = []byte(" autoplay")
	AttrControls                = []byte(" controls")
	AttrHeight                  = []byte(" height=\"")
	AttrWidth                   = []byte(" width=\"")
	AttrLoop                    = []byte(" loop")
	AttrMuted                   = []byte(" muted")
	AttrPoster                  = []byte(" poster=\"")
	AttrPreload                 = []byte(" preload=\"")
	AttrLoading                 = []byte(" loading=\"")
	AttrCrossOrigin             = []byte(" crossorigin=\"")
	AttrControlsList            = []byte(" controlslist=\"")
	AttrDisablePictureInPicture = []byte(" disablepictureinpicture")
	AttrDisableRemotePlayback   = []byte(" disableremoteplayback")
	AttrPlaysInline             = []byte(" playsinline")
)
