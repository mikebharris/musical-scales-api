# Musical Scale Generation API

This simple AWS Lambda function depends on my music theory and practice Go module
at https://github.com/mikebharris/music. All the computation is done in that module; this project is just a thin wrapper
to expose it as a web service.

The service outputs **just** and **tempered** musical scales in **JSON** or the **Scala** file format (see https://www.huygens-fokker.org/scala/). The Scala file can be imported into Scala itself and various other music programs.

**API documentation:** https://mikebharris.github.io/musical-scales-api/

## Examples

<details>
 <summary><code>GET</code> <code><b>/scaleLength={scaleLength}</b></code> <code>(returns fret positions for just intonation and scale length of {scaleLength}</code></summary>

##### Parameters

> | name           | type     | data type | default | description                                                                           |
> |----------------|----------|-----------|---------|---------------------------------------------------------------------------------------|
> | `tuningSystem` | required | string    |         | Tuning system to employ                                                               |
> | `mode`         | optional | string    | ionian  | Which diatonic scale to use for Ptolemy's Intense Diatonic (tuningSystem = 'ptolemy') |
> | `limit`        | optional | int       | 5       | Limit for just intonation (tuningSystem = 'justFromRatios')                           |
> | `divisions`    | optional | int       | 31      | Number of divisions of the octave for equal temperament (tuningSystem = 'edo')        |
> | `format`       | optional | string    | json    | Set to "scala" to change the output format to Scala                                   |

##### Values for `tuningSystem`

> | value                 | type     | description                                                                         |
> |-----------------------|----------|-------------------------------------------------------------------------------------|
> | `justFromRatios`      | just     | 5-limit Just Intonation derived from pure ratios                                    |
> | `pythagorean`         | just     | Pythagorean 3-limit just tuning                                                     |
> | `pythagorean5`        | just     | 5-limit Just Intonation derived from tweaking Pythagorean scale by a syntonic comma |
> | `ptolemy`             | just     | Ptolemy's Intense Diatonic tuning                                                   |
> | `saz`                 | just     | Turkish Saz tuning                                                                  |
> | `edo`                 | tempered | Equal Temperament (Equal Divisions of the Octave)                                   |
> | `meantone`            | tempered | Quarter-Comma Meantone                                                              |
> | `extendedMeantone`    | tempered | Extended Quarter-Comma Meantone                                                     |
> | `bachWellTemperament` | tempered | Bach's Well Temperament (as decoded by Bradley Lehman)                              |

##### Values for `mode`

> | value        | description     |
> |--------------|-----------------|
> | `lydian`     | Lydian mode     |
> | `ionian`     | Ionian mode     |
> | `mixolydian` | Mixolydian mode |
> | `dorian`     | Dorian mode     |
> | `aeloian`    | Aeolian mode    |
> | `phrygian`   | Phrygian mode   |
> | `locrian`    | Locrian mode    |

##### Values for `limit`

> | value   | description                                                          |
> |---------|----------------------------------------------------------------------|
> | uint    | A positive integer (usually a prime; common are 3, 5, 7, 11, 13, 17) |

##### Values for `divisions`

> | value | description                                                           |
> |-------|-----------------------------------------------------------------------|
> | uint  | Any positive integer (default = 12; common are 19, 23, 31, 53, 54, 55 |

##### Values for `format`

> | value   | description                                     |
> |---------|-------------------------------------------------|
> | `json`  | Return JSON (the default)                       |
> | `scala` | Return a Scala file (content-type = text/plain) |

##### Responses

> | http code | content-type       | response                                 |
> |-----------|--------------------|------------------------------------------|
> | `200`     | `application/json` | JSON object                              |
> | `200`     | `text/plain`       | Scala file                               |
> | `422`     | `application/json` | `{"code":"422","message":"Bad Request"}` |

##### Example cURL

**Return Ptolemy's Intense Diatonic tuning in Ionian mode in JSON:**

> ```shell
>  curl -X GET -H "Content-Type: application/json" https://someawsgeneratedlambdaid.lambda-url.us-east-1.on.aws/?tuningSystem=ptolemy
> ```

````json
  {
  "system": "Ptolemy Intense Diatonic",
  "description": "Ptolemy's 5-limit intense diatonic scale in Ionian mode.",
  "intervals": [
    "1:1",
    "9:8",
    "5:4",
    "4:3",
    "3:2",
    "5:3",
    "15:8",
    "2:1"
  ]
}
````

**Return Saz scale as a Scala file:** 

> ```shell
>  curl -X GET -H "Content-Type: application/json" https://someawsgeneratedlambdaid.lambda-url.us-east-1.on.aws/?tuningSystem=saz&format=scala
> ```

````text
! saz.scl
! generated using github.com/mikebharris/music/scala
!
Saz scale using Turkish Saz tuning ratios.
17
18/17
12/11
9/8
81/68
27/22
81/64
4/3
24/17
16/11
3/2
27/17
18/11
27/16
16/9
32/17
64/33
2/1
````

</details>

## Building and provisioning

To build this project, copy the
tool https://github.com/mikebharris/aws-deployment-pipeline-orchestration-helper-tool/blob/main/pipeline.go the project
at https://github.com/mikebharris/aws-deployment-pipeline-orchestration-helper-tool into the top-level directory, and
do:

```shell
go mod tidy
go run pipeline.go --help
```

An example Terraform build and deploy command line:

```shell
go run pipeline.go --stage=build
AWS_ACCESS_KEY_ID=???? AWS_SECRET_ACCESS_KEY=???? go run pipeline.go --account-number=123456789012 --app-name=musical-scales-api --environment=prod --region=us-east-1 --stage=apply --confirm=true
````

Refer to the documentation in that project for more details on how to use the deployment helper tool.
