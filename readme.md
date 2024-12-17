# Worklog Management #

A small command line tool to manage worklogs from multiple users, intended to allow a team to manage incident response.  Convention-based file and directory usage, with toml-based frontmatter.  Designed to get a user up and running with a worklog very quickly to avoid cutting into incident response time.

This does not handle synchronization or encryption.  At its core it's just a document template manager.

Usage:

    # create a new incident
    ftgworklog incident

    # documentation
    ftgworklog --help

