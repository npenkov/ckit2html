# CodeKit 3 to HTML Util

The project aims to help those, who have written projects in CodeKit 3, but want to process those files on Linux/CICD Pipeline/Windows.
CodeKit 3 is and excellent tool and we recommend to go and buy it.

The goal here is to have alternative for different platforms.

**DISCLAIMER**: I am not developer of CodeKit, nor I am planning to support any requests in that area.

## Features and limitrations

If you have a look at [Languages: Kit](https://codekitapp.com/help/kit/) - those are the features that this utility tries to fovide in a single binary, that runs on all platforms.

### Current limitations

 * Supports only variables defined like:
    
    ```html
    <!-- $my-variable:MY-VALUE -->
    ```

    or multiline

    ```html
    <!-- 
    $my-variable:
    MY-VALUE 
    -->
    ```

  * Supports only inclusion of variables in the following form:
  
    ```html
    <!-- $my-variable -->
    ```

  * Supports only includes/imports using the following variations

    ```html
    <!-- @import relative_path/file1.kit -->
    ```

    ```html
    <!-- @include relative_path/file1.kit -->
    ```

    ```html
    <!-- @import "relative_path/file1.kit"  'relative_path/file2.kit' -->
    ```

  * DOES NOT Support `@compile` expressions
  * DOES NOT Support `nil` variables
  * Destination directories have to be precreated. The utility does not take care to create new directories and will fail in writing the output files if the directory does not exist

## Installation

Download the latest release from [ckit2html releases](https://github.com/npenkov/ckit2html/releases) that matches you operating system (Windows, MacOSX-darwin, Linux).

Extract the archive and place the utility into you path (suggested is also renaming it also just to `ckit2html`)


## Usage

```sh
Usage of ckit2html:
  -in string
        Input folder (default ".")
  -out string
        Output folder (default ".")
  -v    Set to verbose output
```

Example (transforming `src/` folder contining `.kit` files to folder `dist/` - where the generated HTMLs will be): 

```sh
ckit2html -in src -out dist
```

NOTE: Imports of of files like ` @import my_kit_file.kit ` work also for `_my_kit_file.kit` as this is the CodeKit behavior. Files that start with `_` will not be transformed to HTML files.

## Alternatives

[The Kit Compiler](https://github.com/bdkjones/Kit/)

## License

Copyright (c) 2020 Nick Penkov. All rights reserved. Use of this source code is governed by a MIT-style license that can be found in the [LICENSE](LICENSE) file.