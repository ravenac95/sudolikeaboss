sudolikeaboss - Now you can too!
================================


.. image:: http://api.flattr.com/button/flattr-badge-large.png
    :alt: Flattr this git repo!
    :target: https://flattr.com/submit/auto?user_id=ravenac95&url=http://sudolikeaboss.com&title=sudolikeaboss&description=sudo+like+a+boss+-+now+you+can+too&language=en_GB&tags=sudolikeaboss,passwords,1password&category=software
|
|
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

*This is by far the easiest method, and the one I recommend most.*

::

    $ brew tap ravenac95/sudolikeaboss
    $ brew install sudolikeaboss


Install from source
*******************

Assuming that you have Go installed and you know how to use it's associated
tools...

::
    
    $ go get github.com/ravenac95/sudolikeaboss

The ``sudolikeaboss`` binary should now be in ``$GOPATH/bin/sudolikeaboss``


Install from zip
****************

Download one of the following zips:

- amd64: http://dl.bintray.com/ravenac95/sudolikeaboss/sudolikeaboss_0.2.1_darwin_amd64.zip

.. warning::
    At this time I'm not sure if the 386 version works. In theory it should,
    but I don't have access to a 32-bit machine to test this.

Then, unzip the file and copy it to the desired location for installation (I
suggest ``/usr/local/bin/sudolikeaboss``).

This entire workflow, would look like this::

    $ mkdir sudolikeaboss
    $ cd sudolikeaboss
    $ wget http://dl.bintray.com/ravenac95/sudolikeaboss/sudolikeaboss_0.2.1_darwin_amd64.zip
    $ unzip sudolikeaboss_0.1.0_darwin_amd64.zip
    $ cp sudolikeaboss /usr/local/bin/sudolikeaboss


.. _configure-iterm:

Configure `iterm2`_ to use ``sudolikeaboss``
********************************************

After installing ``sudolikeaboss``, you still need to configure `iterm2`_. This
is fairly simple. Just watch this gif!

.. image:: https://raw.githubusercontent.com/ravenac95/readme-images/master/sudolikeaboss/configuration.gif

.. _onepass5:

Configuring 1Password5 to work with sudolikeaboss
*************************************************

If you're using 1Password5, or you run into this screen:

.. image:: https://raw.githubusercontent.com/ravenac95/readme-images/master/sudolikeaboss/cannot-fill-item-error-popup.png

This causes a problem for ``sudolikeaboss`` as it isn't a "trusted browser" per 
se. In order to fix this issue, you need to change some preferences on your 
1Password installation. Open up 1password's preferences and find the 
``Advanced`` settings tab. Then make sure to uncheck the option 
``Verify browser code signature``. After doing that, ``sudolikeaboss`` 
should work... like a boss. For the visual learners here's a screenshot:

.. image:: https://cloud.githubusercontent.com/assets/889219/6270365/a69a0726-b816-11e4-9b96-558ddeb00378.png


Getting passwords into `1password`_
-----------------------------------

To get `1password`_ to play ball, just make sure that any passwords you set use
``sudolikeaboss://local`` as the website on the 1password UI. Watch this
example:

.. image:: https://raw.githubusercontent.com/ravenac95/readme-images/master/sudolikeaboss/add-password.gif


Potential Plans for the future!
-------------------------------

These are some ideas I have for the future. This isn't an exhaustive list, and,
more importantly, I make no guarantees on whether or not I can or will get to
any of these.

- Ability to save passwords directly from the command line. Of any of these
  plans, this is probably the most feasible. Again, no promises, but I
  personally want this feature too
- ``tmux`` support. So for those of you that don't use iterm2 I may be able to
  create a different kind of plugin that can work with this.
- linux support? This is a big question mark. If I can get tmux support to
  work, then presumably doing something similar for linux wouldn't be
  impossible. However, the other hard part of this is that linux doesn't
  currently have a GUI for 1password, but I actually have plans to attempt to
  create a gui using some already built tools


Gotchas/Known Issues
--------------------

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

Unfortunately, I have no current plans to do this on Windows. This is mostly
because I wouldn't know where to start. At the moment this software is pretty
dependent on somethings like iterm2 and 1password. As my expertise is in
Linux/Unix environments and not in Windows, I'm not exactly sure what
tools/workflow someone in that camp would use. If you'd like to help out in
this arena, I would be more than happy to give it all a shot.


What's that weird icon on the top-right of the iterm2 window?
*************************************************************

That's just an icon that indicates that an iterm2 `coprocess`_ is running. It
will disappear eventually, as ``sudolikeaboss`` times out after 30 seconds when
waiting for user input.

.. _coprocess: https://iterm2.com/coprocesses.html#/section/home


Do you have this "undocumented API" documented somewhere?
*********************************************************

Not yet, but it will happen soon, hopefully.


Contributing/Developing
-----------------------

I would love help on this! This is actually my first Go project. I'm normally a
Python guy, but decided to take this idea and make it a Go project (which has
been great fun). Any suggestions on how to make this more idiomatic and more
awesome are absolutely welcome.
