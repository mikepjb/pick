# Pick

Fuzzy file finder like selecta in around 150 lines of Go

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

## TODO

- Case insensitive search (on lowercase input)
- Ctrl J/K up/down selection

## Contributing

Feel free to make a pull request.

## License
Pick is released under the [MIT License](https://opensource.org/licenses/MIT)
