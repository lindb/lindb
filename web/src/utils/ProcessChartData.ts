import { getChartColor, toRGBA } from './util'
import { getOptions } from '../config/chartConfig'
import { ChartDatasets, ResultSet, UnitEnum } from '../model/Metric'

/**
 * Generate Line Chart data and options
 * @param {ResultSet} resultSet
 * @param {UnitEnum} unit Current chart Y-axes unit
 */
export function LineChart(resultSet: ResultSet | null, unit?: UnitEnum) {
  if (!resultSet || !resultSet.result) {
    return {}
  }

  const { result: { interval, groups, startTime } } = resultSet

  if (!groups || groups.length === 0 || !startTime || !interval) {
    return {}
  }

  // const times = [ ...Array(pointCount) ].join('0').split('').map((_, idx) => startTime + idx * interval)
  // const labels = times.map(time => moment(time).format('HH:mm'))
  const datasets: ChartDatasets[] = []
  let colorIdx = 0

  groups.forEach(item => {
    const { group, fields } = item

    if (!group || !fields) {
      return
    }

    const groupName = Object.keys(group).map(key => group[ key ]).join('/')

    for (let key of Object.keys(fields)) {
      const bgColor = getChartColor(colorIdx++)

      const fill = false
      const borderColor = bgColor
      const backgroundColor = bgColor
      const label = groupName ? groupName : key
      const pointBackgroundColor = toRGBA(bgColor, 0.25)

      const data = fields[ key ].map((value, index) => ({
        x: new Date(startTime + index * interval),
        y: value ? Math.floor(value * 1000) / 1000 : 0,
      }))

      datasets.push({ label, data, fill, backgroundColor, borderColor, pointBackgroundColor })
    }
  })

  const plugins: any[] = [] // Line Plugins

  return {
    data: { datasets },
    options: getOptions(unit),
    plugins,
  }
}