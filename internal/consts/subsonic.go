package consts

const (
	Xmlns                = "http://subsonic.org/restapi"
	SubsonicVersion      = "1.16.1"
	SubsonicMajorVersion = 1
	SubsonicMinorVersion = 16
)

var (
	SubsonicValidBitRates = []string{
		"0",
		"32",
		"40",
		"48",
		"56",
		"64",
		"80",
		"96",
		"112",
		"128",
		"160",
		"192",
		"224",
		"256",
		"320",
	}

	SubsonicValidFileFormats = []string{
		"mp3",
		"flac",
		"wav",
	}

	SubsonicUserRoles = []string{
		"scrobblingEnabled",
		"ldapAuthenticated",
		"adminRole",
		"settingsRole",
		"streamRole",
		"jukeboxRole",
		"downloadRole",
		"uploadRole",
		"playlistRole",
		"coverArtRole",
		"commentRole",
		"podcastRole",
		"shareRole",
		"videoConversionRole",
	}

	SubsonicErrorMessages = map[string]string{
		"0":  "A generic error.",
		"10": "Required parameter is missing.",
		"20": "Incompatible Subsonic REST protocol version. Client must upgrade.",
		"30": "Incompatible Subsonic REST protocol version. Server must upgrade.",
		"40": "Wrong username or password.",
		"41": "Token authentication not supported for LDAP users.",
		"50": "User is not authorized for the given operation.",
		"60": "The trial period for the Subsonic server is over. Please upgrade to Subsonic Premium. Visit subsonic.org for details.",
		"70": "The requested data was not found.",
	}
)
