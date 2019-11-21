import { getOptions } from '../config/chartConfig';
import { ChartDatasets, ResultSet, UnitEnum } from '../model/Metric';
import { getChartColor, toRGBA } from './Util';

/**
 * Generate Line Chart data and options
 * @param {ResultSet} resultSet
 * @param {UnitEnum} unit Current chart Y-axes unit
 */
export function LineChart(resultSet: ResultSet | null, unit?: UnitEnum) {
  if (!resultSet) {
    return {}
  }

  const { series, startTime, endTime, interval } = resultSet

  if (!series || series.length === 0) {
    return {}
  }

  // const times = [ ...Array(pointCount) ].join('0').split('').map((_, idx) => startTime + idx * interval)
  // const labels = times.map(time => moment(time).format('HH:mm'))
  const datasets: ChartDatasets[] = []
  let colorIdx = 0

  series.forEach(item => {
    const { tags, fields } = item

    if (!fields) {
      return
    }

    const groupName = JSON.stringify(tags)

    for (let key of Object.keys(fields)) {
      const bgColor = getChartColor(colorIdx++)

      const fill = false
      const borderColor = bgColor
      const backgroundColor = bgColor
      const label = groupName ? groupName : key
      const pointBackgroundColor = toRGBA(bgColor, 0.25)

      let data: any[] = []
      const points: { [timestamp: string]: number } = fields[key]
      let i = 0;
      let timestamp = startTime! + i * interval!
      for (; timestamp <= endTime!; i++) {
        console.log("i", i)
        const value = points[`${timestamp}`];
        data.push({
          x: new Date(timestamp),
          y: value ? Math.floor(value * 1000) / 1000 : 0,
        })
        i++
        timestamp = startTime! + i * interval!
      }
      datasets.push({ label, data, fill, backgroundColor, borderColor, pointBackgroundColor })
    }
  })
  const plugins: any[] = [] // Line Plugins
  const options = getOptions(unit)
  return {
    type: "line",
    data: { datasets },
    options: options,
    plugins,
  }
}