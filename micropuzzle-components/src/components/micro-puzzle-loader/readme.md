# micro-puzzle-loader



<!-- Auto Generated Below -->


## Properties

| Property             | Attribute            | Description                                      | Type     | Default     |
| -------------------- | -------------------- | ------------------------------------------------ | -------- | ----------- |
| `fallbacks`          | `fallbacks`          | The count of Fallbacks by initial loading        | `number` | `undefined` |
| `pagesstr`           | `pagesstr`           | Declarations of the pages                        | `string` | `undefined` |
| `streamingurl`       | `streamingurl`       | The URL where to connect the streamingconnection | `string` | `undefined` |
| `streamregistername` | `streamregistername` | The UUID given by the Micropuzzle SSI side       | `string` | `undefined` |


## Events

| Event         | Description                                                              | Type                                              |
| ------------- | ------------------------------------------------------------------------ | ------------------------------------------------- |
| `new-content` | This is a internal micropuzzle event. If new Content is there to show it | `CustomEvent<{ content: string; name: string; }>` |


----------------------------------------------

*Built with [StencilJS](https://stenciljs.com/)*
