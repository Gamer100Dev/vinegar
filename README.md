# <img src="splash/vinegar.png" alt="Vinegarlogo"> Vinegar BSD
[![Pipeline Status][pipeline_img    ]][pipeline    ]
[![Version        ][version_img     ]][version     ]
[![Flathub        ][flathub_img     ]][flathub     ]
[![Report Card    ][goreportcard_img]][goreportcard]
[![Discord Server ][discord_img     ]][discord     ]
[![Matrix Room    ][matrix_img      ]][matrix      ]

An open-source, minimal, configurable, fast bootstrapper for running Roblox on BSD and Linux.

[pipeline]:     https://github.com/vinegarhq/vinegar/actions/workflows/go.yml
[pipeline_img]: https://img.shields.io/github/actions/workflow/status/vinegarhq/vinegar/go.yml?style=flat-square&label=build%20%26%20tests
[version]:     https://github.com/vinegarhq/vinegar/releases/latest
[version_img]: https://img.shields.io/github/v/release/vinegarhq/vinegar?style=flat-square&display_name=tag
[flathub]:     https://flathub.org/apps/details/org.vinegarhq.Vinegar
[flathub_img]: https://img.shields.io/flathub/downloads/org.vinegarhq.Vinegar?style=flat-square
[goreportcard]:     https://goreportcard.com/report/github.com/vinegarhq/vinegar
[goreportcard_img]: https://goreportcard.com/badge/github.com/vinegarhq/vinegar?style=flat-square
[discord]:     https://discord.gg/dzdzZ6Pps2
[discord_img]: https://img.shields.io/discord/1069506340973707304?style=flat-square&label=discord
[matrix]:      https://matrix.to/#/#vinegarhq:matrix.org
[matrix_img]:  https://img.shields.io/matrix/vinegarhq:matrix.org?style=flat-square&label=matrix

# Features
+ Multiple instances of Roblox open simultaneously
+ Automatic Wineprefix killer when Roblox has quit
+ Logging for both Vinegar, Wine and Roblox (Currently Bugged)
+ Modifications of Roblox via the Overlay directory, overwriting Roblox's files; such as re-adding the old death sound
+ Automatic DXVK Installer and uninstaller
+ Fast Multi-threaded installation and extraction of Roblox
+ Automatic removal of outdated cached packages and versions of Roblox
+ Custom execution of wine program within wineprefix
+ Set different environment variables and FFlags for both Player and Studio, with Global to override
+ Force a specific version of Roblox to be deployed
+ Custom launcher specified to be used when launching Roblox
+ Wine Root feature to set a specific wine installation path
+ Sanitization of environment
+ Browser launch via MIME
+ Splash window during setup, with error dialog support

## BSD Features:

+ Roblox and Studio launches
+ Workes with any wine version (without any wow64 errors)
+ Detection of discord is present ( you may add the protocols to your client to suit it)

## To come features

+ Rewrite Gamemode (linux) to Gameboost available at: github.com/thindil/gameboost
+ Fixes for Roblox's player and studio; patches will be added to this repo regarding wine
+ Fix for the client
+ Login portal integration patches for studio
  
# See Also
+ [Discord server](https://discord.gg/dzdzZ6Pps2)
+ [Matrix room](https://matrix.to/#/#vinegarhq:matrix.org)
+ [Documentation](https://vinegarhq.github.io)
+ [Roblox-Studio-Mod-Manager](https://github.com/MaximumADHD/Roblox-Studio-Mod-Manager)
+ [Bloxstrap](https://github.com/pizzaboxer/bloxstrap)

# Acknowledgements
+ Credits to
  + [pizzaboxer](https://github.com/pizzaboxer)
  + [MaximumADHD](https://github.com/MaximumADHD)
+ Logo modified with Katie. Base glass SVG made by the [Twemoji team](https://twemoji.twitter.com/), Licensed under [CC-BY 4.0](https://creativecommons.org/licenses/by/4.0/).
+ Katie usage authorized by [karliflux](https://karliflux.neocities.org)
