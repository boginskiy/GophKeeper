package pkg

import (
	"fmt"
	"os"
)

var Keeper = `
                                                                                                                  
KKKKKKKKK    KKKKKKK                                                                                                    
K:::::::K    K:::::K                                                                                                    
K:::::::K    K:::::K                                                                                                    
K:::::::K   K::::::K                                                                                                    
KK::::::K  K:::::KKK    eeeeeeeeeeee        eeeeeeeeeeee    ppppp   ppppppppp       eeeeeeeeeeee    rrrrr   rrrrrrrrr   
  K:::::K K:::::K     ee::::::::::::ee    ee::::::::::::ee  p::::ppp:::::::::p    ee::::::::::::ee  r::::rrr:::::::::r  
  K::::::K:::::K     e::::::eeeee:::::ee e::::::eeeee:::::eep:::::::::::::::::p  e::::::eeeee:::::eer:::::::::::::::::r 
  K:::::::::::K     e::::::e     e:::::ee::::::e     e:::::epp::::::ppppp::::::pe::::::e     e:::::err::::::rrrrr::::::r
  K:::::::::::K     e:::::::eeeee::::::ee:::::::eeeee::::::e p:::::p     p:::::pe:::::::eeeee::::::e r:::::r     r:::::r
  K::::::K:::::K    e:::::::::::::::::e e:::::::::::::::::e  p:::::p     p:::::pe:::::::::::::::::e  r:::::r     rrrrrrr
  K:::::K K:::::K   e::::::eeeeeeeeeee  e::::::eeeeeeeeeee   p:::::p     p:::::pe::::::eeeeeeeeeee   r:::::r            
KK::::::K  K:::::KKKe:::::::e           e:::::::e            p:::::p    p::::::pe:::::::e            r:::::r            
K:::::::K   K::::::Ke::::::::e          e::::::::e           p:::::ppppp:::::::pe::::::::e           r:::::r            
K:::::::K    K:::::K e::::::::eeeeeeee   e::::::::eeeeeeee   p::::::::::::::::p  e::::::::eeeeeeee   r:::::r            
K:::::::K    K:::::K  ee:::::::::::::e    ee:::::::::::::e   p::::::::::::::pp    ee:::::::::::::e   r:::::r            
KKKKKKKKK    KKKKKKK    eeeeeeeeeeeeee      eeeeeeeeeeeeee   p::::::pppppppp        eeeeeeeeeeeeee   rrrrrrr            
                                                             p:::::p                                                    
                                                             p:::::p                                                    
                                                            p:::::::p                                                   
                                                            p:::::::p                                                   
                                                            p:::::::p                                                   
                                                            ppppppppp                                                   

`

func PrintInfo(version, date, commit string) {
	if version == "" {
		version = "N/A"
	}
	if date == "" {
		date = "N/A"
	}
	if commit == "" {
		commit = "N/A"
	}
	fmt.Println(Keeper)
	fmt.Fprintf(os.Stdout, "Build version: %s\n", version)
	fmt.Fprintf(os.Stdout, "Build date:    %s\n", date)
	fmt.Fprintf(os.Stdout, "Build commit:  %s\n", commit)
}
