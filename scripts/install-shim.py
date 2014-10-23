import os
import sys
import shutil

STARTING_COMMENT = "//START SUDOLIKEABOSS SHIM"
END_COMMENT = "//END SUDOLIKEABOSS SHIM"
EXTENSION_NAME = "aomjjhallfgjeglblehebfpbcfeobpgk"

PATH_TO_EXTENSION = os.path.join(
    os.path.expanduser('~/Library/Application Support/Google/Chrome/Default/Extensions'),
    EXTENSION_NAME
)

CURRENT_DIR = os.path.dirname(os.path.abspath(__file__))

CANNOT_FIND_EXTENSION_ERROR_MESSAGE = """
Error: Cannot find correct Chrome Extension

In order for the sudolikeaboss fix for 1password to work please install the
Chrome Extension 4.2.5 or greater.
"""

UNSUPPORTED_EXTENSION_ERROR_MESSAGE = """
Error: Chrome Extension found but version is incompatible

The file that this shim is looking cannot be found. You seem to be using a
greater version of this extension which is unsupported
"""


def cannot_find_extension():
    print CANNOT_FIND_EXTENSION_ERROR_MESSAGE
    sys.exit(1)


def unsupported_extension():
    print UNSUPPORTED_EXTENSION_ERROR_MESSAGE
    sys.exit(1)


def main():
    print "Installing the shim for sudolikeaboss 1Password 5 support"

    if not os.path.isdir(PATH_TO_EXTENSION):
        cannot_find_extension()

    try:
        version_dir_name = os.listdir(PATH_TO_EXTENSION)[0]
    except:
        cannot_find_extension()

    # Ensure the version numbers are ready
    version_number = version_dir_name.split('_')[0]

    version_list = map(lambda a: int(a), version_number.split('.'))

    # if the version is less than 4.2.0 then i dunno if this works so don't let
    # it do anything for now
    if version_list < [4, 2, 0]:
        cannot_find_extension()

    background_js_file_path = os.path.join(PATH_TO_EXTENSION, version_dir_name,
                                           'code', 'global.min.js')

    if not os.path.isfile(background_js_file_path):
        unsupported_extension()

    print "Creating a backup of the background script as ./background.js"
    shutil.copyfile(background_js_file_path, 'background.js')

    background_js_file = open(background_js_file_path)

    background_js_file_lines = []

    collect = True

    # Check for already existing file
    for line in background_js_file.readlines():
        if line.startswith(STARTING_COMMENT):
            collect = False
        if collect:
            background_js_file_lines.append(line.strip())
        if line.startswith(END_COMMENT):
            collect = True
    background_js_file.close()

    # Append the javascript to this file
    shim_js_file = open(os.path.join(CURRENT_DIR, 'chrome-extension-shim.js'))

    shim_js_str = shim_js_file.read()

    # Add the shim
    background_js_file_lines.append(STARTING_COMMENT)
    background_js_file_lines.append(shim_js_str)
    background_js_file_lines.append(END_COMMENT)

    # Write extension background file with the shim
    background_js_file = open(background_js_file_path, 'w')
    background_js_file.write('\n'.join(background_js_file_lines))
    background_js_file.close()


if __name__ == '__main__':
    main()
