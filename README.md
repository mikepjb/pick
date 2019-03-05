# Pick

Fuzzy file finder

Current progress: works very basically on the command line with single level
pipes:  
`ls | pick`

However it does not deal well with secondary pipes: `cat $(ls | pick)`
