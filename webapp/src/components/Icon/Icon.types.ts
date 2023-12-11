export type IconName = 'user' | 'message' | 'connectAccount' | 'warning' | 'close' | 'globe' | 'msTeams' | 'link' | 'noChannels' | 'lock' | 'mattermost' | 'tick'

export type IconProps = {
    iconName: IconName;
    height?: number;
    width?: number;
    className?: string;
}
