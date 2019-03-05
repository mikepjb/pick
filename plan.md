# How does the result get output in selecta?

- L55 it is puts (writes to stdout afaik)

# How is the usable output for vim seperated from what the user sees while
operating selecta?

- Possibly the Screen.with_screen on L45
  - Screen.with_screen calls TTY.with_tty
    - Gary has given us a helpful note starting L796 which reads:
      Selecta reads data from stdin and writes it to stdout, so we can't draw
      UI and receive keystrokes through them. Fortunately, all modern
      Unix-likes provide /dev/tty, which IO.console gives us.

# What is TTY
  - Again from Gary's text L795, it looks as though TTY is a virtual terminal,
    which sits on top of an IO.console => it's invoke with TTY.new(console_file)
  - console_file = IO.console

# Where does user input come from?
  - After newing up the TTY, get_available_input is a function that calls .getc
    (get character) on the console_file

# How can I open an interactive subprogram from within a go program?
https://www.reddit.com/r/golang/comments/2nd4pq/how_can_i_open_an_interactive_subprogram_from/

cmd := exec.Command("vim", "filename")
cmd.Run()

cmd.Stdout = os.Stdout
cmd.Stdin = os.Stdin
cmd.Std = os.Stderr
cmd.Run()

Interesting but not quite what I need..

# How do I consume stdout as stdin without writing back to stdout?
  Gary's L826 in Selecta uses IO.pipe

# How do I print to the console without writing to Stdout
  - I think /dev/tty is the answer, Gary uses it in selecta.

# What are /dev/tty and /dev/stdout?

Not sure but stdout seems to be both for display and passing info to another
process. tty is just for display.

I *think* they are just files.. this is unix after all..

# there is a terminal package in golang x/crypto...
The first thing I tried didn't work, didn't pursue it very far as I believe
/dev/tty is a file and should be easy to manipulate without extra wrapping code
that I did not write.

# Can I open /dev/tty as a file and write to it?

Yes you can but you won't see the output on your terminal (even after flushing)

However `echo 1 > /dev/tty` prints 1 on the current terminal.

I think /dev/tty is the current file for the output of your terminal.
weirdly tailing it produces no output.. probably because it's pointing to the
current "/dev/tty" for the terminal you're running in currently.

stdout as a side effect prints here too?

# How do I read from tty during pipe?
