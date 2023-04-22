/*
 * @Author: fyfishie
 * @Date: 2023-04-22:11
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-04-22:14
 * @@email: fyfishie@outlook.com
 * @Description: :)
 */
package lib

type RDNSResItem struct {
	IP      string   `json:"ip"`
	Domains []string `json:"domains"`
}
