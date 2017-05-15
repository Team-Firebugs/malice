package plugin

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/server/httputils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/streamformatter"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

func parseHeaders(headers http.Header) (map[string][]string, *types.AuthConfig) {

	metaHeaders := map[string][]string{}
	for k, v := range headers {
		if strings.HasPrefix(k, "X-Meta-") {
			metaHeaders[k] = v
		}
	}

	// Get X-Registry-Auth
	authEncoded := headers.Get("X-Registry-Auth")
	authConfig := &types.AuthConfig{}
	if authEncoded != "" {
		authJSON := base64.NewDecoder(base64.URLEncoding, strings.NewReader(authEncoded))
		if err := json.NewDecoder(authJSON).Decode(authConfig); err != nil {
			authConfig = &types.AuthConfig{}
		}
	}

	return metaHeaders, authConfig
}

// parseRemoteRef parses the remote reference into a reference.Named
// returning the tag associated with the reference. In the case the
// given reference string includes both digest and tag, the returned
// reference will have the digest without the tag, but the tag will
// be returned.
func parseRemoteRef(remote string) (reference.Named, string, error) {
	// Parse remote reference, supporting remotes with name and tag
	remoteRef, err := reference.ParseNormalizedNamed(remote)
	if err != nil {
		return nil, "", err
	}

	type canonicalWithTag interface {
		reference.Canonical
		Tag() string
	}

	if canonical, ok := remoteRef.(canonicalWithTag); ok {
		remoteRef, err = reference.WithDigest(reference.TrimNamed(remoteRef), canonical.Digest())
		if err != nil {
			return nil, "", err
		}
		return remoteRef, canonical.Tag(), nil
	}

	remoteRef = reference.TagNameOnly(remoteRef)

	return remoteRef, "", nil
}

func (pr *pluginRouter) getPrivileges(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := httputils.ParseForm(r); err != nil {
		return err
	}

	metaHeaders, authConfig := parseHeaders(r.Header)

	ref, _, err := parseRemoteRef(r.FormValue("remote"))
	if err != nil {
		return err
	}

	privileges, err := pr.backend.Privileges(ctx, ref, metaHeaders, authConfig)
	if err != nil {
		return err
	}
	return httputils.WriteJSON(w, http.StatusOK, privileges)
}

func (pr *pluginRouter) upgradePlugin(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := httputils.ParseForm(r); err != nil {
		return errors.Wrap(err, "failed to parse form")
	}

	var privileges types.PluginPrivileges
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&privileges); err != nil {
		return errors.Wrap(err, "failed to parse privileges")
	}
	if dec.More() {
		return errors.New("invalid privileges")
	}

	metaHeaders, authConfig := parseHeaders(r.Header)
	ref, tag, err := parseRemoteRef(r.FormValue("remote"))
	if err != nil {
		return err
	}

	name, err := getName(ref, tag, vars["name"])
	if err != nil {
		return err
	}
	w.Header().Set("Docker-Plugin-Name", name)

	w.Header().Set("Content-Type", "application/json")
	output := ioutils.NewWriteFlusher(w)

	if err := pr.backend.Upgrade(ctx, ref, name, metaHeaders, authConfig, privileges, output); err != nil {
		if !output.Flushed() {
			return err
		}
		output.Write(streamformatter.FormatError(err))
	}

	return nil
}
