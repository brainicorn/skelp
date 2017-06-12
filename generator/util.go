package generator

import (
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

const (
	TIDTypeRepo    = "REPO"
	TIDTypeFile    = "FILE"
	TIDTypeAlias   = "ALIAS"
	TIDTypeUnknown = "UNKNOWN"
)

var (
	scpLikeUrlRegExp   = regexp.MustCompile("^(ssh://)?(?:(?P<user>[^@]+)(?:@))?(?P<host>[^:|/]+)(?::|/)?(?P<path>.+)$")
	isHttpSchemeRegExp = regexp.MustCompile("^(http|https)://")
	isFileSchemeRegExp = regexp.MustCompile("^(file)://")
)

func TypeForTemplateID(templateID string) string {
	if IsRepoURL(templateID) {
		return TIDTypeRepo
	}

	if IsFilePath(templateID) {
		return TIDTypeFile
	}

	return TIDTypeAlias
}

func AuthMethodForURL(url string) transport.AuthMethod {
	var am transport.AuthMethod

	if IsSSH(url) {
		// 1:scheme, 2:user, 3:host, 4:path
		m := scpLikeUrlRegExp.FindStringSubmatch(url)
		a, err := ssh.NewSSHAgentAuth(m[2])

		if err == nil {
			am = a
		}
	}

	return am
}

func FilepathFromURL(u string) (string, error) {
	var err error
	var fpath string
	var pu *url.URL

	if !IsRepoURL(u) {
		return fpath, fmt.Errorf("invalid url %s", u)
	}

	if isHttpSchemeRegExp.MatchString(u) {
		pu, err = url.Parse(u)

		if err == nil {
			fpath = pu.Hostname() + pu.Path
		}
	} else if IsSSH(u) {
		// 1:scheme, 2:user, 3:host, 4:path
		m := scpLikeUrlRegExp.FindStringSubmatch(u)
		fpath = m[3] + "/" + m[4]
	}

	if strings.HasSuffix(fpath, ".git") {
		fpath = strings.TrimSuffix(fpath, ".git")
	}

	fpath = filepath.FromSlash(fpath)
	return fpath, err
}

func IsRepoURL(templateID string) bool {
	return isHttpSchemeRegExp.MatchString(templateID) || scpLikeUrlRegExp.MatchString(templateID)
}

func IsSSH(templateID string) bool {
	return !isHttpSchemeRegExp.MatchString(templateID) && scpLikeUrlRegExp.MatchString(templateID)
}

func IsFilePath(templateID string) bool {
	return isFileSchemeRegExp.MatchString(templateID) || (!IsRepoURL(templateID) && strings.Contains(filepath.ToSlash(templateID), "/"))
}

func IsAlias(templateID string) bool {
	return !strings.Contains(filepath.ToSlash(templateID), "/")
}
