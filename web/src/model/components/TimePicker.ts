import moment from 'moment-es6'

export class TimePickerModel {
  private title: string
  private from: string
  private to: string
  private time: {
    from: number
    to: number
  }

  constructor(title: string, from: string, to: string) {
    this.title = title
    this.from = from
    this.to = to

    this.time = {
      from: moment(this.from).valueOf(),
      to: moment(this.to).valueOf(),
    }
  }
}

/**
 * parse string like `now()-30s`、`now()-1h`
 * @param {string} time
 */
function parseTime(time: string) {
  // ...
}