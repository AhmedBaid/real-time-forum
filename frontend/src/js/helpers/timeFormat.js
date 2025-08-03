export function timeFormat(time) {
    let timef = ""
    let postDate = new Date(time).getTime()
    let now = new Date()
    now.setHours(now.getHours() )
    let def = (now - postDate) / 1000
    if (def < 60) {
        timef = "just now"
    } else if (def < 3600) {
        let m = Math.floor(def / 60)
        timef = `${m} minutes ago `
    } else if (def < 86400) {
        let h = Math.floor(def / 3600)
        timef = `${h} hours ago `
    } else {
        let d = Math.floor(def / 86400)
        timef = `${d} days ago `
    }

    return timef
}


