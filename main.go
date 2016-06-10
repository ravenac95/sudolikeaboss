package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>
int
StartApp(void) {
	[NSAutoreleasePool new];
	[NSApplication sharedApplication];
	[NSApp setActivationPolicy:NSApplicationActivationPolicyProhibited];
	[NSApp run];
	return 0;
}
*/
import "C"

import (
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()

	app.Name = "sudolikeaboss"
	app.Version = "0.2.1"
	app.Usage = "use 1password from the terminal with ease"
	app.Action = func(c *cli.Context) {
		go runSudolikeaboss()
		C.StartApp()
	}

	app.Run(os.Args)
}
