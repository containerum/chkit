// +build !dev
// +build !mock

package cmd

// DEBUG -- true if dev version
const DEBUG = false

// MOCK -- true if dev with mock api version
const MOCK = false

// API_ADDR -- gateway API addr
var API_ADDR = "https://192.168.88.200:8082"
