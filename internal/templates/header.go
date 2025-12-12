package templates

import "fmt"

// Header prints an ASCII header.
func Header() {
	fmt.Println(`
    __                        ___ 
   / /___   ______ ___  _____/ (_)
  / //_/ | / / __ '__ \/ ___/ / / 
 / ,<  | |/ / / / / / / /__/ / /  
/_/|_| |___/_/ /_/ /_/\___/_/_/
    `)
}
