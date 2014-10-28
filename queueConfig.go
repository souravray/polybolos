/*
* @Author: souravray
* @Date:   2014-10-12 01:31:40
* @Last Modified by:   souravray
* @Last Modified time: 2014-10-28 01:20:09
 */

package polybolos

// type queueConfig struct {
// 	Name string
// 	// Number of tries/leases after which the task fails permanently and is deleted.
// 	// If AgeLimit is also set, both limits must be exceeded for the task to fail permanently.
// 	RetryLimit int32

// 	// Maximum time allowed since the task's first try before the task fails permanently and is deleted (only for push tasks).
// 	// If RetryLimit is also set, both limits must be exceeded for the task to fail permanently.
// 	AgeLimit time.Duration

// 	// Minimum time between successive tries (only for push tasks).
// 	MinBackoff time.Duration

// 	// Maximum time between successive tries (only for push tasks).
// 	MaxBackoff time.Duration

// 	// Maximum number of times to double the interval between successive tries before the intervals increase linearly (only for push tasks).
// 	MaxDoublings int32

// 	// If MaxDoublings is zero, set ApplyZeroMaxDoublings to true to override the default non-zero value.
// 	// Otherwise a zero MaxDoublings is ignored and the default is used.
// 	ApplyZeroMaxDoublings bool
// }
