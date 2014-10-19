sudolikeaboss - Now you can too!
================================

.. image:: https://raw.githubusercontent.com/ravenac95/readme-images/master/sudolikeaboss/demo.gif

Pretty neat, eh?


What's happening here?
----------------------

``sudolikeaboss`` is a simple application that aims to make your life as a dev,
ops, or just a random person who likes to ssh and sudo into boxes much, much
easier by allowing you to access your `1password` passwords on the terminal.
All you need is `iterm2`_, `1password`_, a mac, and a dream.

.. _iterm2: http://iterm2.com/
.. _1password: https://agilebits.com/onepassword


Benefits
--------

- Better security through use of longer, more difficult to guess passwords
- Better security now that you can have a different password for every server
  if you'd like
- Greater convenience when accessing passwords on the terminal


So is this only for sudo passwords?
-----------------------------------

No! You can use this for tons of things! Like...

- `dm-crypt`_ passwords on external boxes
- `gpg`_ passwords to use on the terminal

.. _dm-crypt: https://code.google.com/p/cryptsetup/wiki/DMCrypt
.. _gpg: https://www.gnupg.org/


Ok! I want it. How do I install this thing?!
--------------------------------------------

I tried to make installation as simple as possible. So here's the quickest path
to awesomeness.

Install with homebrew
*********************

::

    $ brew tap ravenac95/sudolikeaboss
    $ brew install sudolikeaboss


Install from source
*******************

*I may eventually add some docs on how to install without homebrew, but really
you should just use homebrew.*


Configure `iterm2`_ to use ``sudolikeaboss``
********************************************

After installing ``sudolikeaboss``, you still need to configure `iterm2`_. This
is fairly simple. Just watch this gif!

.. image:: https://raw.githubusercontent.com/ravenac95/readme-images/master/sudolikeaboss/configuration.gif


Getting passwords into `1password`_
-----------------------------------

To get `1password`_ to play ball, just make sure that any passwords you set use
``sudolikeaboss://local`` as the website on the 1password UI. Watch this
example:

.. image:: https://raw.githubusercontent.com/ravenac95/readme-images/master/sudolikeaboss/add-password.gif


Gotchas
-------

Here are just some questions or gotchas that I figured people would run into or
have.


Why is the 1password popup not where I'm typing?
************************************************

The way the popup works is by finding your mouse cursor. I'd like to improve
this, but since I'm using 1password's undocumented API this is how it will be
right now.


I don't use 1password
*********************

Are you serious?! If you're on a mac and you have passwords, you should be
using 1password. With that said, I would love to support additional password
managers as the project grows. 


I use linux
***********

Sorry :( I don't have anything for you yet. Maybe you can help me with that :)


I use Windows
*************

Sorry? *But seriously, I have no current plans to do this on Windows.*


What's that weird icon on the top-right of the iterm2 window?
*************************************************************

That's just an icon that indicates that an iterm2 `coprocess`_ is running. It
will disappear eventually, as ``sudolikeaboss`` times out after 30 seconds when
waiting for user input.

.. _coprocess: https://iterm2.com/coprocesses.html#/section/home


Contributing/Developing
-----------------------

I would love help on this! This is actually my first Go project. I'm normally a
Python guy, but decided to take this idea and make it a Go project (which has
been great fun). Any suggestions on how to make this more idiomatic and more
awesome are absolutely welcome.
