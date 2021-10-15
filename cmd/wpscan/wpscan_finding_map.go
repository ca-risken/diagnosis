package main

const (
	typeHeaders            = "headers"
	typeSearchReplaceDB2   = "search_replace_db2"
	typeReadme             = "readme"
	typeDebugLog           = "debug_log"
	typeBackupDB           = "backup_db"
	typeFullPathDisclosure = "full_path_disclosure"
)

var wpscanFindingMap = map[string]wpscanFindingInformation{
	typeHeaders:            {Score: 1.0, Description: "Software version found by Headers"},
	typeSearchReplaceDB2:   {Score: 8.0, Description: ""},
	typeReadme:             {Score: 6.0, Description: ""},
	typeDebugLog:           {Score: 8.0, Description: ""},
	typeBackupDB:           {Score: 6.0, Description: ""},
	typeFullPathDisclosure: {Score: 6.0, Description: ""},
}

type wpscanFindingInformation struct {
	Score       float32
	Description string
}
