package constant

import "sync"

var (
	// 线程安全的map
	DefaultIMG = "https://images.669pic.com/element_min_new_pic/91/81/13/10/9a2397491c92eb8674c90a62ec41b0ca.png"
	DBLock     sync.RWMutex
)

var PushTemplate = `**{_name}**({_size} GB)

豆瓣(-) | IMDB(-) | TMDB(-)

#{_from}

频道:@炸鸡巴牛逼`

var BoxTemplate = `
{name}
D: {dlspeed}/s  U: {upspeed}/s
D: {downloaded} U: {uploaded}
R: {availability}
{state}
{time}
`

var BoxSystemInfo = `
	D: {dl_info_speed}/s  U: {up_info_speed}/s
	D: {dl_info_data} U: {up_info_data}
	limit : {dl_rate_limit}/s  {up_rate_limit}/s
	{box_name}
	CPU:{CPU}
	RAM:{RAM}
	IO:{IO}
`
