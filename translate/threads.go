package main

import "net/url"

func getMainThread(ts, channel string) string {
	perm, err := getPermalink(ts, channel)
	if err != nil {
		return ""
	}

	if !perm.OK {
		return ""
	}

	u, err := url.Parse(perm.Permalink)
	if err != nil {
		return ""
	}

	tsparam := u.Query().Get("thread_ts")
	if tsparam == "" || tsparam != ts {
		return tsparam
	}

	return ""
}
