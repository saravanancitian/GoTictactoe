package com.samaya.gotictactoe;


import java.util.Date;


public class ScoreSettings implements java.io.Serializable{
    public static final long SEC_IN_MILLIS = 1000;
    private int totalPlayed;
    private int totalWin;
    private int totalTied;

    long date1;
    long date2;

    public boolean isSettingsSound() {
        return settingsSound;
    }

    public boolean isSettingShowTimer() {
        return settingShowTimer;
    }

    public void setSettingsSound(boolean settingsSound) {
        this.settingsSound = settingsSound;
    }

    public void setSettingShowTimer(boolean settingShowTimer) {
        this.settingShowTimer = settingShowTimer;
    }

    boolean settingsSound = true;

    boolean settingShowTimer = true;

    public long getDate1() {
        return date1;
    }

    public long getDate2() {
        return date2;
    }

    public long getDate3() {
        return date3;
    }

    public long getTopPlayedTime1() {
        return topPlayedTime1;
    }

    public long getTopPlayedTime2() {
        return topPlayedTime2;
    }

    public long getTopPlayedTime3() {
        return topPlayedTime3;
    }

    long date3;
    long topPlayedTime1;
    long topPlayedTime2;
    long topPlayedTime3;

    private static ScoreSettings Instance;
    private ScoreSettings(){

        reset();
    }


    public static ScoreSettings getInstance() {
        if(Instance == null){
            Instance = new ScoreSettings();
        }
        return Instance;
    }

    public int getTotalPlayed() {
        return totalPlayed;
    }

    public int getTotalWin() {
        return totalWin;
    }

    public int getTotalTied() {
        return totalTied;
    }

    private void addTopTime(long curPlayedTime){

        if(topPlayedTime1 == 0){
            topPlayedTime1 = curPlayedTime;
            date1 = System.currentTimeMillis();
        } else if(curPlayedTime <= topPlayedTime1) {
            topPlayedTime3 = topPlayedTime2;
            date3 = date2;
            topPlayedTime2 = topPlayedTime1;
            date2 = date1;
            topPlayedTime1 = curPlayedTime;
            date1 = System.currentTimeMillis();
        } else if (curPlayedTime <= topPlayedTime2){
            topPlayedTime3 = topPlayedTime2;
            date3 = date2;
            topPlayedTime2 = curPlayedTime;
            date2 = System.currentTimeMillis();
        } else if (curPlayedTime <= topPlayedTime3){
            topPlayedTime3 = curPlayedTime;
            date3 = System.currentTimeMillis();
        }
    }

    public  void addPlayed(int winner, long duration){
        this.totalPlayed++;
        if(winner == MainActivity.GAME_TIED) {
            this.totalTied++;
        } else if (winner == MainActivity.HUMAN_PLAYER){
            this.totalWin++;
            addTopTime(duration);
        }

    }

    public void reset(){
        totalPlayed = 0;
        totalWin = 0;
        totalTied = 0;
        date1 = 0;
        date2 = 0;
        date3 = 0;
        topPlayedTime1 = topPlayedTime2 = topPlayedTime3 = 0;
    }


}
