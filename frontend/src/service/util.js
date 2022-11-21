export default function toDateTime(unix) {
    let datetime = new Date(unix * 1000)
    return datetime.toLocaleDateString("id-ID", { weekday: 'long' }) + ', ' + datetime.getUTCDate() + ' ' + datetime.toLocaleDateString("id-ID", { month: 'long' }) + ' ' + datetime.getFullYear() + ' ' + datetime.getHours() + ':' + datetime.getMinutes() + ':' + datetime.getSeconds()
}