// 类型定义（用于JavaScript代码提示）

/**
 * @typedef {Object} Point
 * @property {number} x
 * @property {number} y
 * @property {number} [pressure]
 */

/**
 * @typedef {Object} DrawStyle
 * @property {string} color
 * @property {number} width
 * @property {number} opacity
 * @property {string} [fill]
 * @property {number} [fillOpacity]
 */

/**
 * @typedef {Object} BaseElement
 * @property {string} id
 * @property {'pen'|'line'|'rect'|'circle'|'ellipse'|'triangle'|'arrow'|'text'|'eraser'|'image'} type
 * @property {Point[]} points
 * @property {DrawStyle} style
 * @property {string} userId
 * @property {number} timestamp
 * @property {string} [text]
 * @property {number} [fontSize]
 * @property {string} [fontFamily]
 * @property {string} [imageData]
 * @property {number} [x]
 * @property {number} [y]
 * @property {number} [width]
 * @property {number} [height]
 */

/**
 * @typedef {Object} User
 * @property {string} id
 * @property {string} name
 * @property {string} avatar
 * @property {string} color
 * @property {Point} [cursor]
 * @property {boolean} connected
 */

/**
 * @typedef {Object} Board
 * @property {string} id
 * @property {string} name
 * @property {string} createdAt
 * @property {string} updatedAt
 * @property {string} background
 * @property {string} [thumbnail]
 */

/**
 * @typedef {Object} HistoryRecord
 * @property {string} id
 * @property {string} boardId
 * @property {string} action
 * @property {string} userId
 * @property {string} userName
 * @property {string} timestamp
 * @property {string} [snapshot]
 */

/**
 * @typedef {Object} WsMessage
 * @property {'draw'|'cursor'|'user_join'|'user_leave'|'users'|'clear'|'undo'|'redo'|'sync'} type
 * @property {*} data
 * @property {string} [userId]
 * @property {number} [timestamp]
 */

export const Point = {}
export const DrawStyle = {}
export const BaseElement = {}
export const User = {}
export const Board = {}
export const HistoryRecord = {}
export const WsMessage = {}
