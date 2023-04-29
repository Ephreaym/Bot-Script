package basicmap_old

/*

void MapInitialize()
{
    WarBot = Object("WarBot");
    WizBot = Object("WizBot");
    WarBotCorpse = Object("WarBotCorpse");
    WizBotCorpse = Object("WizBotCorpse");
    OutOfGameWar = Waypoint("OutOfGameWar");
    OutOfGameWiz = Waypoint("OutOfGameWiz");
    WarSound = Waypoint("WarSound");
    WizSound = Waypoint("WizSound");
    WarCryCooldown = 0;
    EyeOfTheWolfCooldown = 0;
    BerserkerChargeCooldown = 0;
    GlobalCooldown = 0;
    RespawnCooldownDelay = 0;
    SecondTimer(2,Respawn);
    SecondTimer(2,RespawnWiz);
}
void CheckCurrentHealth()
{
    if (CurrentHealth(WarBot) < 1001)
    {
        ChangeScore(GetCaller(),1);
        DeadBotWar();
    }
    if (CurrentHealth(WizBot) < 1001)
    {
        ChangeScore(GetCaller(),1);
        DeadBotWiz();
    }
}
void DeadBotWar()
{
    MoveWaypoint(WarSound,GetObjectX(WarBot),GetObjectY(WarBot));
    AudioEvent("NPCDie",WarSound);
    MoveObject(WarBotCorpse,GetObjectX(WarBot),GetObjectY(WarBot));
    MoveObject(WarBot,GetWaypointX(OutOfGameWar),GetWaypointY(OutOfGameWar));
    SecondTimer(2,Respawn);
}
void DeadBotWiz()
{
    MoveWaypoint(WizSound,GetObjectX(WizBot),GetObjectY(WizBot));
    AudioEvent("NPCDie",WizSound);
    MoveObject(WizBotCorpse,GetObjectX(WizBot),GetObjectY(WizBot));
    MoveObject(WizBot,GetWaypointX(OutOfGameWiz),GetWaypointY(OutOfGameWiz));
    SecondTimer(2,RespawnWiz);
}
void Respawn()
{
    RestoreHealth(WarBot, MaxHealth(WarBot) - CurrentHealth(WarBot));
    WarCryCooldown = 0;
    EyeOfTheWolfCooldown = 0;
    BerserkerChargeCooldown = 0;
    MoveWaypoint(WarSound,GetObjectX(WarBotCorpse),GetObjectY(WarBotCorpse));
    AudioEvent("BlinkCast",WarSound);
    BotSpawn = Object("BotSpawn" + IntToString(Random(1,14)));
    MoveObject(WarBot,GetObjectX(BotSpawn),GetObjectY(BotSpawn));
    MoveObject(WarBotCorpse,GetWaypointX(OutOfGameWar),GetWaypointY(OutOfGameWar));
    Enchant(WarBot,"ENCHANT_INVULNERABLE",5.0);
    RespawnCooldownDelay = 1;
    SecondTimer(10,RespawnCooldownDelayReset);
}
void RespawnWiz()
{
    RestoreHealth(WizBot, MaxHealth(WizBot) - CurrentHealth(WizBot));
    MoveWaypoint(WizSound,GetObjectX(WizBotCorpse),GetObjectY(WizBotCorpse));
    AudioEvent("BlinkCast",WizSound);
    BotSpawn = Object("BotSpawn" + IntToString(Random(1,14)));
    MoveObject(WizBot,GetObjectX(BotSpawn),GetObjectY(BotSpawn));
    MoveObject(WizBotCorpse,GetWaypointX(OutOfGameWiz),GetWaypointY(OutOfGameWiz));
    Enchant(WizBot,"ENCHANT_INVULNERABLE",5.0);
}
void WarCry()
{
    if (MaxHealth(other) == 150)
    {
    }
    else
    {
        if ((WarCryCooldown == 0) && (GlobalCooldown == 0))
        {
            MoveWaypoint(WarSound,GetObjectX(WarBot),GetObjectY(WarBot));
            AudioEvent("WarcryInvoke",WarSound);
            PauseObject(WarBot,45);
            Enchant(WarBot,"ENCHANT_HELD",1.0);
            Enchant(other,"ENCHANT_ANTI_MAGIC",3.0);
            EnchantOff(self,"ENCHANT_SHOCK");
            EnchantOff(self,"ENCHANT_INVULNERABLE");
            WarCryCooldown = 1;
            GlobalCooldown = 1;
            SecondTimer(10,WarCryCooldownReset);
            SecondTimer(1,GlobalCooldownReset);
        }
        else
        {
        }
    }
}
void WarCryCooldownReset()
{
    if (RespawnCooldownDelay == 0)
    {
        WarCryCooldown = 0;
    }
    else
    {
    }
}
void EyeOfTheWolf()
{
    Wander(WarBot);
    if (EyeOfTheWolfCooldown == 0)
    {
    Enchant(self,"ENCHANT_INFRAVISION",10.0);
    EyeOfTheWolfCooldown = 1;
    SecondTimer(20,EyeOfTheWolfCooldownReset);
    }
    else
    {
    }
}
void EyeOfTheWolfCooldownReset()
{
    if (RespawnCooldownDelay == 0)
    {
        EyeOfTheWolfCooldown = 0;
    }
    else
    {
    }
}

// ------------------------------------- BERSERKERCHARGE SCRIPT ---------------------------------------------- //


float UnitRatioX(int unit, int target, float size)
{
    return (GetObjectX(unit) - GetObjectX(target)) * size / Distance(GetObjectX(unit), GetObjectY(unit), GetObjectX(target), GetObjectY(target));
}

float UnitRatioY(int unit, int target, float size)
{
    return (GetObjectY(unit) - GetObjectY(target)) * size / Distance(GetObjectX(unit), GetObjectY(unit), GetObjectX(target), GetObjectY(target));
}
float ToFloat(int x)
{
    StopScript(x);
}
int ToInt(float x)
{
    StopScript(x);
}
void WarBotDetectEnemy()
{
    if ((BerserkerChargeCooldown == 0) && (GlobalCooldown == 0))
    {
        int rnd = Random(0, 2);

        if (!rnd || (rnd == 1))
        {
            BerserkerChargeCooldown = 1;
            GlobalCooldown = 1;
            SecondTimer(1,GlobalCooldownReset);
            BerserkerInRange(GetTrigger(), GetCaller(), 10);
        }
    }
    else
    {
        if (WarCryCooldown == 0)
        {
            WarCry();
        }
        else
        {
        }
    }
}
int CheckUnitFrontSight(int unit, float dtX, float dtY)
{
    MoveWaypoint(1, GetObjectX(unit) + dtX, GetObjectY(unit) + dtY);
    int temp = CreateObject("InvisibleLightBlueHigh", 1);
    int res = IsVisibleTo(unit, temp);

    Delete(temp);
    return res;
}
void pointedByPlr()
{
}
void BerserkerInRange(int owner, int target, int wait)
{
    int unit;

    if (CurrentHealth(owner) && CurrentHealth(target))
    {
        if (!HasEnchant(owner, "ENCHANT_ETHEREAL"))
        {
            Enchant(owner, "ENCHANT_ETHEREAL", 0.0);
            MoveWaypoint(1, GetObjectX(owner), GetObjectY(owner));
            unit = CreateObject("InvisibleLightBlueHigh", 1);
            MoveWaypoint(1, GetObjectX(unit), GetObjectY(unit));
            CreateObject("InvisibleLightBlueHigh", 1);
            Raise(unit, ToFloat(owner));
            Raise(unit + 1, ToFloat(target));
            LookWithAngle(unit, wait);
            FrameTimerWithArg(1, unit, BerserkerWaitStrike);
        }
    }
}

void BerserkerWaitStrike(int ptr)
{
    int count = GetDirection(ptr), owner = ToInt(GetObjectZ(ptr)), target = ToInt(GetObjectZ(ptr + 1));

    while (1)
    {
        if (IsObjectOn(ptr) && CurrentHealth(owner) && CurrentHealth(target) && IsObjectOn(owner))
        {
            if (count)
            {
                if (IsVisibleTo(owner, target) && Distance(GetObjectX(owner), GetObjectY(owner), GetObjectX(target), GetObjectY(target)) < 400.0)
                    BerserkerCharge(owner, target);
                else
                {
                    LookWithAngle(ptr, count - 1);
                    FrameTimerWithArg(6, ptr, BerserkerWaitStrike);
                    break;
                }
            }
        }
        if (CurrentHealth(owner))
            EnchantOff(owner, "ENCHANT_ETHEREAL");
        if (IsObjectOn(ptr))
        {
            Delete(ptr);
            Delete(ptr + 1);
        }
        break;
    }
}

void BerserkerCharge(int owner, int target)
{
    int unit;
    if (CurrentHealth(owner) && CurrentHealth(target))
    {
        EnchantOff(owner,"ENCHANT_INVULNERABLE");
        MoveWaypoint(2, GetObjectX(owner), GetObjectY(owner));
        AudioEvent("BerserkerChargeInvoke", 2);
        MoveWaypoint(1, GetObjectX(owner), GetObjectY(owner));
        unit = CreateObject("InvisibleLightBlueHigh", 1);
        MoveWaypoint(1, GetObjectX(unit), GetObjectY(unit));
        Raise(CreateObject("InvisibleLightBlueHigh", 1), UnitRatioX(target, owner, 23.0));
        Raise(CreateObject("InvisibleLightBlueHigh", 1), UnitRatioY(target, owner, 23.0));
        CreateObject("InvisibleLightBlueHigh", 1);
        Raise(unit + 3, ToFloat(owner));
        LookWithAngle(GetLastItem(owner), 0);
        SetCallback(owner, 9, BerserkerTouched);
        Raise(unit, ToFloat(target));
        LookAtObject(unit + 1, target);
        FrameTimerWithArg(1, unit, BerserkerLoop);
    }
}

void BerserkerLoop(int ptr)
{
    int owner = ToInt(GetObjectZ(ptr + 3)), count = GetDirection(ptr);

    if (CurrentHealth(owner) && count < 60 && IsObjectOn(ptr) && IsObjectOn(owner))
    {
        if (CheckUnitFrontSight(owner, GetObjectZ(ptr + 1) * 1.5, GetObjectZ(ptr + 2) * 1.5) && !GetDirection(GetLastItem(owner)))
        {
            MoveObject(owner, GetObjectX(owner) + GetObjectZ(ptr + 1), GetObjectY(owner) + GetObjectZ(ptr + 2));
            LookWithAngle(owner, GetDirection(ptr + 1));
            Walk(owner, GetObjectX(owner), GetObjectY(owner));
        }
        else
            LookWithAngle(ptr, 100);
        FrameTimerWithArg(1, ptr, BerserkerLoop);
    }
    else
    {
        SetCallback(owner, 9, NullCollide);
        Delete(ptr);
        Delete(ptr + 1);
        Delete(ptr + 2);
        Delete(ptr + 3);
    }
}

void BerserkerTouched()
{
    if (IsObjectOn(self))
    {
        while (1)
        {
            if (!GetCaller() || (HasClass(other, "IMMOBILE") && !HasClass(other, "DOOR") && !HasClass(other, "TRIGGER")) && !HasClass(other, "DANGEROUS"))
            {
                MoveWaypoint(2, GetObjectX(self), GetObjectY(self));
                AudioEvent("FleshHitStone", 2);

                Enchant(self, "ENCHANT_HELD", 2.0);
            }
            else if (CurrentHealth(other))
            {
                if (IsAttackedBy(self, other))
                {
                    MoveWaypoint(2, GetObjectX(self), GetObjectY(self));
                    AudioEvent("FleshHitFlesh", 2);
                    Damage(other, self, 100, 2);
                }
                else
                break;
            }
            else
            break;
            LookWithAngle(GetLastItem(self), 1);
            break;
        }
    }
    Wander(WarBot);
    SecondTimer(10,BerserkerChargeCooldownReset);
}
void NullCollide()
{
    return;
}
void BerserkerChargeCooldownReset()
{
    if (RespawnCooldownDelay == 0)
    {
        BerserkerChargeCooldown = 0;
    }
    else
    {
    }
}
void GlobalCooldownReset()
{
    GlobalCooldown = 0;
}
void RespawnCooldownDelayReset()
{
    RespawnCooldownDelay = 0;
}

int BotSpawn;
int WarBot;
int WizBot;
int WizBotCorpse;
int WarBotCorpse;
int OutOfGameWar;
int OutOfGameWiz;
int WarSound;
int WizSound;
int WarCryCooldown;
int EyeOfTheWolfCooldown;
int BerserkerChargeCooldown;
int GlobalCooldown;
int RespawnCooldownDelay;

*/
