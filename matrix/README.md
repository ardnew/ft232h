# matrix
This program constructs the Travis CI configuration (YAML) file. 

Travis currently does not support conditional jobs based on target architecture, so I've resorted to creating my own explicit build matrix from a simple Go struct which defines all possible build+test target combinations.

This program takes no arguments and prints the complete YAML content to stdout.
