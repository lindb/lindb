import { UnitEnum } from '../model/Metric'
const convert = require('convert-units')

export class DataFormatter {
  public static formatter(point: number, unit: UnitEnum) {
    switch (unit) {
      case UnitEnum.Nanoseconds:
        return this.transformNanoSeconds(point);
      case UnitEnum.Milliseconds:
        return this.transformMilliseconds(point)
      case UnitEnum.Seconds:
        return this.transformSeconds(point)
      case UnitEnum.Bytes:
        return this.transformBytes(point)
      case UnitEnum.Percent:
        return this.transformPercent(point)
      default:
        return this.transformNone(point)
    }
  }

  public static transformNanoSeconds(input: number, decimals?: number): string {
    if (!input) {
      return "0";
    }
    const best = convert(input).from("ns").toBest();
    const value = convert(input).from("ns").to(best.unit);
    if (decimals !== undefined) {
      return value.toFixed(decimals) + best.unit;
    } else {
      return Math.floor(value * 100) / 100 + best.unit;
    }
  }

  public static transformMilliseconds(input: number): string {
    if (input > 24 * 3600 * 1000) {
      return (input / (24 * 3600 * 1000)).toFixed(2) + ' days'
    } else if (input > 3600 * 1000) {
      return (input / (3600 * 1000)).toFixed(2) + ' hours'
    } else if (input > 10 * 60 * 1000) {
      return (input / 60000).toFixed(2) + ' min'
    } else if (input > 1000) {
      return (input / 1000).toFixed(2) + ' s'
    } else if (!input) {
      return '0 ms'
    } else {
      return input.toFixed(2) + ' ms'
    }
  }

  public static transformSeconds(input: number): string {
    if (input > 365 * 24 * 3600) {
      return (input / (365 * 24 * 3600)).toFixed(2) + ' years'
    } else if (input > 24 * 3600) {
      return (input / (24 * 3600)).toFixed(2) + ' days'
    } else if (input > 3600) {
      return (input / (3600)).toFixed(2) + ' hours'
    } else if (input > 60) {
      return (input / 60).toFixed(2) + ' minutes'
    } else if (!input) {
      return '0 sec'
    } else {
      return input.toFixed(2) + ' sec'
    }
  }

  public static transformBytes(input: number): string {
    if (input > 1024 * 1024 * 1024 * 1024 * 1024) {
      return (input / (1024 * 1024 * 1024 * 1024 * 1024)).toFixed(2) + ' PB'
    } else if (input > 1024 * 1024 * 1024 * 1024) {
      return (input / (1024 * 1024 * 1024 * 1024)).toFixed(2) + ' TB'
    } else if (input > 1024 * 1024 * 1024) {
      return (input / (1024 * 1024 * 1024)).toFixed(2) + ' GB'
    } else if (input > 1024 * 1024) {
      return (input / (1024 * 1024)).toFixed(2) + ' MB'
    } else if (input > 1024) {
      return (input / 1024).toFixed(2) + ' KB'
    } else if (!input) {
      return '0 Byte'
    } else {
      return input.toString() + ' Byte'
    }
  }

  public static transformPercent(input: number): string {
    if (!input) {
      return '0%'
    } else {
      return input.toFixed(2).toString() + '%'
    }
  }

  public static transformNone(input: number): string {
    if (input > 1000 * 1000 * 1000) {
      return (input / (1000 * 1000 * 1000)).toFixed(2) + 'B'
    } else if (input > 1000 * 1000) {
      return (input / (1000 * 1000)).toFixed(2) + 'M'
    } else if (input > 1000) {
      return (input / 1000).toFixed(2) + 'K'
    } else if (!input) {
      return '0'
    } else {
      return input.toString() + ''
    }
  }
}