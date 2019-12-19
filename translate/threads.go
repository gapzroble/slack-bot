package main

import "net/url"

func getMainThread(ts, channel string, tsChan chan<- string) {
	threadTs := ts
	defer func() {
		tsChan <- threadTs
	}()

	perm, err := getPermalink(ts, channel)
	if err != nil {
		return
	}

	if !perm.OK {
		return
	}

	u, err := url.Parse(perm.Permalink)
	if err != nil {
		return
	}

	tsparam := u.Query().Get("thread_ts")
	if tsparam != "" {
		threadTs = tsparam

	}
}
