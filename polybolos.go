/*
* Πολυβολος
* @Author: souravray
* @Date:   2014-10-11 19:52:00
* @Last Modified by:   souravray
* @Last Modified time: 2015-02-16 02:28:00
 */

// Polybolos is a feature rich embedded Go job-queue,
// for executing longer tasks in the background.
package polybolos

import (
	"github.com/souravray/polybolos/sys"
)

// Initialize the pacakge with required settings
func init() {
	sys.UseMaxCPUs()
}

// Returns current version
func Version() string {
	return "alpha"
}
