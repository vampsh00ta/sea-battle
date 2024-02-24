PROJECTNAME=seabattle

init:
	find . -type  f -exec sed -i '' -e's/seabattle/$(PROJECTNAME)/g' {} + -not -path "Makefile"
