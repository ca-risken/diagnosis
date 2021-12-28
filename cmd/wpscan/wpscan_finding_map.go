package main

const (
	typeVersion                     = "Version"
	typeVersionInsecure             = "Version/Insecure"
	typeHeaders                     = "headers"
	typeSearchReplaceDB2            = "search_replace_db2"
	typeLoginClosed                 = "LoginPage/Closed"
	typeLoginOpened                 = "LoginPage/Opened"
	typeLoginOpenedUserFound        = "LoginPage/Opened/UserFound"
	typeReadme                      = "readme"
	typeDebugLog                    = "debug_log"
	typeBackupDB                    = "backup_db"
	typeFullPathDisclosure          = "full_path_disclosure"
	recommendTypeSearchReplaceDB2   = "SearchReplaceDB"
	recommendTypeReadme             = "readme.html"
	recommendTypeDebugLog           = "debug.log"
	recommendTypeBackupDB           = "backup-db"
	recommendTypeFullPathDisclosure = "FullPathDisclosure"
)

var wpscanFindingMap = map[string]wpscanFindingInformation{
	typeVersion: {Score: 1.0, Description: "WordPress version %v identified"},
	typeVersionInsecure: {Score: 6.0, Description: "WordPress version %v identified (Insecure)", RecommendType: typeVersionInsecure,
		Risk: `WordPress Insecure Version
	- WordPress is not up to date and not secure. 
	- Vulnerabilities may exist, and attacks can cause serious damage.`,
		Recommendation: `Update wordpress.
	- https://wordpress.org/support/article/updating-wordpress/`},
	typeHeaders: {Score: 1.0, Description: "Software version found by Headers"},
	typeSearchReplaceDB2: {Score: 8.0, Description: "", RecommendType: recommendTypeSearchReplaceDB2,
		Risk: `Search Replace DB script found
	- It may be possible to manipulate the database in the Web UI.`,
		Recommendation: `Delete search-replace-db directory.
	- Restrict access to search-replace-db directory.
	- https://github.com/interconnectit/Search-Replace-DB`},
	typeReadme: {Score: 6.0, Description: "", RecommendType: recommendTypeReadme,
		Risk: `Readme.html exists
	- Basic information such as version can be identified, which may provide useful information to the attacker.`,
		Recommendation: `Delete readme.html.`},
	typeDebugLog: {Score: 8.0, Description: "", RecommendType: recommendTypeDebugLog,
		Risk: `debug.log exists
	- An attacker can use the debugging information in this file to conduct further attacks.`,
		Recommendation: `Disable debug mode.
	- Restrict access to debug.log.
	- https://wordpress.org/support/article/debugging-in-wordpress/`},
	typeBackupDB: {Score: 6.0, Description: "", RecommendType: recommendTypeBackupDB,
		Risk: `backup-db directory exists
	- An Attacker may be able to obtain database backups from outside.`,
		Recommendation: `Delete backup-db directory.
	- Restrict access to backup-db directory.`},
	typeFullPathDisclosure: {Score: 6.0, Description: "", RecommendType: recommendTypeFullPathDisclosure,
		Risk: `The fully qualified path is displayed 
	- The fully qualified path is displayed in PHP error messages or other pages.
	- An attacker can use the the file structure of the web server to conduct further attacks.`,
		Recommendation: `Prevent this information from being displayed to the user.`},
	typeLoginClosed: {Score: 1.0, Description: "WordPress login page is closed."},
	typeLoginOpened: {Score: 8.0, Description: "WordPress login page is opened.", RecommendType: typeLoginOpened,
		Risk: `Login page is opened
	- If weak passwords are used or usernames are identifiable, an attack may gain access to restricted content.`,
		Recommendation: `Restrict access to the admin panel.`},
	typeLoginOpenedUserFound: {Score: 9.0, Description: "WordPress login page is open. And username was found.", RecommendType: typeLoginOpenedUserFound,
		Risk: `Login page is opened
	- If weak passwords are used or usernames are identifiable, an attack may gain access to restricted content.
	- WordPress username has been identified by scan`,
		Recommendation: `Restrict access to the admin panel.`},
}

type wpscanFindingInformation struct {
	Score          float32
	Description    string
	RecommendType  string
	Risk           string
	Recommendation string
}
