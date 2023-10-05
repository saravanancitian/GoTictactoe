package com.samaya.gotictactoe;


import java.util.Date;


public class Score  implements java.io.Serializable{
    private static final long SEC_IN_MILLIS = 1000;
    private int totalPlayed;
    private int totalWin;
    private int totalTied;

    long date1;
    long date2;
    long date3;
    long topPlayedTime1;
    long topPlayedTime2;
    long topPlayedTime3;

    private static Score Instance;
    private Score(){

        reset();
    }


    public static Score getInstance() {
        if(Instance == null){
            Instance = new Score();
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
        if (curPlayedTime >= topPlayedTime1){
            topPlayedTime3 = topPlayedTime2;
            date3 = date2;
            topPlayedTime2 = topPlayedTime1;
            date2 = date1;
            topPlayedTime1 = curPlayedTime;
            date1 = System.currentTimeMillis();
        } else if (curPlayedTime >= topPlayedTime2){
            topPlayedTime3 = topPlayedTime2;
            date3 = date2;
            topPlayedTime2 = curPlayedTime;
            date2 = System.currentTimeMillis();
        } else if(curPlayedTime >= topPlayedTime3){
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

    public String formattedString(String instr){

        Date dt1 = date1 > 0? new Date(date1 ): new Date();
        long dur1 =  topPlayedTime1 > 0? topPlayedTime1/SEC_IN_MILLIS : 0;
        Date dt2 =  date2 > 0? new Date(date2) : new Date();
        long dur2 =  topPlayedTime2 > 0? topPlayedTime2/SEC_IN_MILLIS: 0;
        Date dt3 =  date3 > 0?new Date(date3) : new Date();
        long dur3 =  topPlayedTime3 > 0? topPlayedTime3/SEC_IN_MILLIS: 0;


        return String.format(instr, totalPlayed, totalWin, totalTied
                ,dt1,dt1,dt1 , dur1/60, dur1 % 60
                ,dt2,dt2,dt2 , dur2/60, dur2 % 60
                ,dt3,dt3,dt3 , dur3/60, dur3 % 60);
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
