# Pick

Fuzzy file finder like selecta in around 150 lines of Go

## Design

go through each candidate (links passed to pick):
  - for each character in the search string:
    - if you match, highlight it and continue to the next
      (specifically avoid multiple highlights of the same word OR avoid double
      counting?!?)
    - rank by character match and length between first and last match

## Usage

On the command line:  
`ls | pick`  
`vim $(find . | pick)`

Inside Vim with the following snippet:  
**Note: does not work on Neovim**

```
func! Pick(input_cmd, vim_cmd)
  try
    let selection = system(a:input_cmd . " | pick")
  catch /Vim:Interrupt/
  endtry
  redraw!
  if len(selection) != 0
    exec a:vim_cmd . " " . selection
  endif
endfunc

nnoremap <leader>f :call Pick("find * -type f", ":e")<cr>
```

## Contributing

Feel free to make a pull request.

## License
Pick is released under the [MIT License](https://opensource.org/licenses/MIT)
