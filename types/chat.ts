export const CHAT_EXHIBIT = {
  ProblemList: 'problemList'
} as const;

export type ChatExhibit = typeof CHAT_EXHIBIT[keyof typeof CHAT_EXHIBIT]; 