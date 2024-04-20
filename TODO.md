# TODO

[x] Download & parse html from url
    [x] Find ingredients by simple search: css classname or id containing the substring `ingredient`
        [x] First look for `ul` or `ol` elements
        [x] Then look for all other types of elements
    [x] Find the cooking instructions/method
        [x] First look for `ul` or `ol` elements
        [x] Then look for all other types of elements
    [ ] A more robust way to search for recipe sections
        [ ] Pass to LLM with a prompt to find anything that looks like a recipe section?
    [ ] Timeout & cancellation
[ ] Testing
    [x] Basic integration test
        [x] Asynchronously pull recipes from multiple sources
        [x] Readable struct/json diffs
    [ ] Better integration testing
        [ ] HTML fixtures
        [ ] Option to update fixtures with "network" run
    [ ] Unit testing
        [ ] fincRecipeSection function
            [x] Basic table-driven tests
            [ ] Tests for finding sections via different element types
                [x] Priority element types: ul, ol
                [ ] Generic element type: *
            [ ] Tests against multiple fixtures
    [ ] Add benchmarks
[x] Asynchronous downloading and processing of each page
[ ] HTTP service
    [ ] Webserver/application server
        [ ] Logging middleware
        [ ] Rendering to templates
        [ ] HTMX
        [ ] API endpoints
            [ ] GET /
            [ ] POST /recipe/pull/{url}
    [ ] Restructure project directory and packages
[ ] Infrastrucutre
    [ ] Compiling/build step
    [ ] Containerization & deployment
