package domain

// User is a type where the user attributes are stored
type User struct {
	ID          int64
	Username    string
	AccessToken string
}

// Link defines the structure to the navigation links
type Link struct {
	Title string
	URL   string
}

type Repository struct {
	Name        *string `json:"name,omitempty"`
	FullName    *string `json:"full_name,omitempty"`
	Description *string `json:"description,omitempty"`
	Private     *bool   `json:"private,omitempty"`
	HTMLURL     *string `json:"html_url,omitempty"`
	CloneURL    *string `json:"clone_url,omitempty"`
	SSHURL      *string `json:"ssh_url,omitempty"`
}

type Key struct {
	ID    *int    `json:"id,omitempty"`
	Key   *string `json:"key,omitempty"`
	Title *string `json:"title,omitempty"`
	URL   *string `json:"url,omitempty"`
}
