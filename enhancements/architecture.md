# Gmail-bot Architecture

## Gmail Templates
All templates should have an according cronjob with email content constructed from `gmail-template.yaml`
```yaml
kind: email # Type of template*
cronjob: # Reoccurring time for sending the email*
to: # Dest of the email*
cc: # cc of the email
title: # Title of the email*
body: # Body of the email which should act as templates and include variable and if statements such as https://golang.org/pkg/text/template/*
signature: # Signature of the email
```
> Note: Fields with \* are required. 

# Testing

All testing are part of the GH workflow.

## On PR Commit

 - Basic build, vendor, and unit tests. 
 - Gmail templates should have their according workflow for sending the email. 
