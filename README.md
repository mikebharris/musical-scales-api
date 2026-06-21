# Musical Scale Generation API (AWS Lambda Service)

This simple AWS Lambda function depends on my music theory and practice Go module
at https://github.com/mikebharris/music. All the computation is done in that module; this project is just a thin wrapper
to expose it as a web service.

The service outputs **just** and **tempered** musical scales in **JSON** or the **Scala** file format (see https://www.huygens-fokker.org/scala/). The Scala file can be imported into Scala itself and various other music programs.

**API documentation:** https://mikebharris.github.io/musical-scales-api/

## Examples

<details>
 <summary><code>GET</code> <code><b>/tuningSystem=edo</b></code> <code>(returns scale intervals for 12-EDO scale length of {scaleLength}</code></summary>

##### Parameters

> | name           | type     | data type | default | description                                                                           |
> |----------------|----------|-----------|---------|---------------------------------------------------------------------------------------|
> | `tuningSystem` | required | string    |         | Tuning system to employ                                                               |
> | `mode`         | optional | string    | Ionian  | Which diatonic scale to use for Ptolemy's Intense Diatonic (tuningSystem = 'ptolemy') |
> | `limit`        | optional | uint      | 5       | Limit for just intonation (tuningSystem = 'justFromRatios')                           |
> | `divisions`    | optional | uint      | 31      | Number of divisions of the octave for equal temperament (tuningSystem = 'edo')        |
> | `partials`     | optional | uint      | 1       | Number of partials of the Harmonic Series to compute (tuningSystem = 'harmonic')      |
> | `format`       | optional | string    | json    | Set to "scala" to change the output format to Scala                                   |

##### Values for `tuningSystem`

> | value                 | type     | description                                                                                              |
> |-----------------------|----------|----------------------------------------------------------------------------------------------------------|
> | `harmonic`            | just     | Ratios generated from the Harmonic Series, quantised into semitones to produce the justest 12-note scale |
> | `justFromRatios`      | just     | 5-limit Just Intonation derived from pure ratios, quantised into semitones with justest ratio chosen     |
> | `partch`              | just     | Harry Partch's 43-tone "Genesis" scale build on 11-limit tuning                                          |
> | `pythagorean`         | just     | Pythagorean 3-limit just tuning                                                                          |
> | `pythagorean5`        | just     | 5-limit Just Intonation derived from tweaking Pythagorean scale by a syntonic comma                      |
> | `ptolemy`             | just     | Ptolemy's Intense Diatonic tuning                                                                        |
> | `saz`                 | just     | Turkish Saz tuning                                                                                       |
> | `edo`                 | tempered | Equal Temperament (Equal Divisions of the Octave)                                                        |
> | `meantone`            | tempered | Quarter-Comma Meantone                                                                                   |
> | `extendedMeantone`    | tempered | Extended Quarter-Comma Meantone                                                                          |
> | `bachWellTemperament` | tempered | Bach's Well Temperament (as decoded by Bradley Lehman)                                                   |

##### Values for `mode`

> | value        | description     |
> |--------------|-----------------|
> | `Lydian`     | Lydian mode     |
> | `Ionian`     | Ionian mode     |
> | `Mixolydian` | Mixolydian mode |
> | `Dorian`     | Dorian mode     |
> | `Eeloian`    | Aeolian mode    |
> | `Phrygian`   | Phrygian mode   |
> | `Locrian`    | Locrian mode    |

##### Values for `limit`

> | value   | description                                                          |
> |---------|----------------------------------------------------------------------|
> | uint    | A positive integer (usually a prime; common are 3, 5, 7, 11, 13, 17) |

##### Values for `divisions`

> | value | description                                                            |
> |-------|------------------------------------------------------------------------|
> | uint  | Any positive integer (default = 12; common are 19, 23, 31, 53, 54, 55) |

##### Values for `partials`

> | value | description                                       |
> |-------|---------------------------------------------------|
> | uint  | Any positive integer (a good example would be 45) |

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
