i want to build a keyboard tester in golang
here are the basic requirements:

- the binary should take a flag in as argument with the type of keyboard: 40%, 60%, 65%, 75%, 80%, 96%, 100%
- this binary will just assume macos, nothing else
- if there are no arguments, it should just display a help

keys testing:
- it should overwrite any settings if the keys might cause some actions on a computer
- a lowercase/uppercase key press should be registered as the same
- i would like to test if the keys can be press and held too 
- any key presses that are not recognized should also be recorded -- but it's not expected to be held
- a combination of keys could be used for any appropriate features (e.g. hotkeys, quit, rest, etc.)

UI:
- the ui should find a simple way to display all the expected keys
- on the side, it should also be able to display all the unknown keys 
- you can use bubbletea to build a simple testing
- a simple instruction and help should be given on the UI
- give a modern look to it, utilizing the features that bubbletea can provide

Code requirements:
- functions should be kept short and easily readable, it should generally not exceed 25 lines, with a hard limit of 30
- files should be aptly name and reasonably contain the right functions
- consts should be set aside in a file; each keyboard config should have an array of key presses it wants to check
- the key mapping should also be put aside in another file that clearly outlines the keypresses that i'm checking for
- the logic should be clear, avoid complex conditionals 

## Clarifications

These clarifications take precedence wherever they conflict with the original requirements above.

### Keyboard capture
- Temporarily intercept keyboard events during active testing to suppress shortcuts and computer actions.
- Use the required macOS permissions; leave saved system preferences unchanged.
- Let the user select one physical keyboard to test. Only input from that keyboard contributes to test results.

### Keyboard layouts
- Use one documented, conventional US ANSI reference layout for each supported keyboard size.
- Distinguish left and right modifier keys.

### Key testing
- Track key presses only.
- Hold testing and key repetition checks are no longer required.
- Do not provide a separate held-key colour or hold completion state.

### Results
- Show test results and unknown keys in the UI.
- Store results only in RAM for the current session.
- Saving reports to disk is not required.
