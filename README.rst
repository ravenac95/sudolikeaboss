sudolikeaboss - Now you can too!
================================

.. image:: https://raw.githubusercontent.com/ravenac95/readme-images/master/sudolikeaboss/demo.gif

Pretty neat, eh? 


What's happening here?
----------------------

``sudolikeaboss`` is a simple application that aims to make your life as a dev,
ops, or just a random person who likes to ssh and sudo into boxes much, much
easier by allowing you to access your `1password` passwords on the terminal.
All you need is `iterm2`_, `1password`_ and a dream.

.. _iterm2: http://iterm2.com/
.. _1password: https://agilebits.com/onepassword


Benefits
--------

- Better security through use of longer, more difficult to guess passwords
- Better security now that you can have a different password for every server
  if you'd like
- Greater convenience when accessing passwords


So is this only for sudo passwords?
-----------------------------------

No! You can use this for tons of things! Like...

- dm-crypt passwords on external boxes
- gpg passwords to use on the terminal


Ok! I want it. How do I install this thing?!
--------------------------------------------

I tried to make installation as simple as possible. So here's the quickest path
to awesomeness.


Install the ``sudolikeaboss``
*****************************

Install with homebrew::

    $ brew tap ravenac95/sudolikeaboss
    $ brew install sudolikeaboss


Configure iterm2 - so you can sudo, like a boss
***********************************************

To setup `iterm2`_ is fairly simple. Just watch this gif!

.. image:: https://raw.githubusercontent.com/ravenac95/readme-images/master/sudolikeaboss/setup.gif


Getting passwords into `1password`_
-----------------------------------

To get `1password`_ to play ball, just make sure that any passwords you set use
``sudolikeaboss://local`` as the website on the 1password UI.


Contributing/Developing
-----------------------

I would love help on this! This is actually my first Go project. I'm normally a
Python guy, but decided to take this idea and make it a Go project (which has
been great fun). Any suggestions on how to make this more idiomatic and more
awesome are absolutely welcome.
