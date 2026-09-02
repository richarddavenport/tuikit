package ui

import "github.com/richarddavenport/tuikit/comp"

// Every region democtl draws, named once.
//
// A capture script and a right-click menu address these by name — `click
// services.row[2]` rather than a coordinate — so the list is worth reading in
// one place. A name that nothing draws is a name the canvas will report as
// missing, which is how a script that has drifted from the interface finds out.
const (
	regHeader comp.Name = "header"
	regRule   comp.Name = "header.rule"
	regFooter comp.Name = "footer"

	regServices    comp.Name = "services"
	regServicesRow comp.Name = "services.row"

	regDetail    comp.Name = "detail"
	regDetailTab comp.Name = "detail.tab"

	regLogs    comp.Name = "logs"
	regLogsRow comp.Name = "logs.row"

	regRun     comp.Name = "run"
	regRunStep comp.Name = "run.step"
	// The meter under the steps, and the bar inside it. Two names because the
	// pixel layer is placed on an OWNER: the picture covers the track and not
	// the brackets or the label around it.
	regRunMeter comp.Name = "run.meter"
	regRunTrack comp.Name = "run.track"

	// regSplit is the column between the panes. It draws nothing — an owned
	// blank is still trimmed from the output — but it is what a drag grabs.
	regSplit comp.Name = "split"

	regConfirm  comp.Name = "confirm"
	regMenu     comp.Name = "menu"
	regMenuItem comp.Name = "menu.item"
)
