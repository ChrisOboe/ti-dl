# ti-dl
A download tool for a music streaming service. 

## Configuration
ti-dl can be either configured through a configuration file or through parameters.

The default path for the config file is
 * Linux: $HOME/.config/ti-dl/config.toml
 * Windows: $APPDATA$\ti-dl\config.toml 
 
 The config file format is
 
 ```
 [User]
 Username=Your Username
 Password=Your Password
 
 [Paths]
 Destination=./${ALBUMARTIST}/${RELEASETYPE}/${RELEASEDATE} - ${RELEASETITLE}/${TRACKNUMBER} - ${TRACKTITLE}
 ```
 
 The following values are possible for paths:
  * ${ALBUMARTIST}: The artist who released the album
  * ${RELEASETYPE}: If its an Album or a Single or whatever
  * ${RELEASEDATE}: The date when it was released
  * ${RELEASETITLE} The title of the release
  * ${TRACKNUMBER}: The tracknumber
  * ${TRACKTITLE}: The tracktitle
  
