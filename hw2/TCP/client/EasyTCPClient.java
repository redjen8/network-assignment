package hw2.TCP.client;

import java.net.*;
import java.nio.charset.StandardCharsets;
import java.util.Scanner;
import java.util.NoSuchElementException;
import java.io.*;

public class EasyTCPClient{

    private static final String SERVER_NAME = "localhost";
    private static final int SERVER_PORT = 22848;

    public static void main(String[] args) {
        Socket socket = null;

        Runtime.getRuntime().addShutdownHook(new Thread(() -> { System.out.println("Bye bye~"); }));

        try {
            socket = new Socket();
            socket.connect(new InetSocketAddress(SERVER_NAME, SERVER_PORT));
            byte[] inputBuffer = null;
            String recvMessage = null;
            OutputStream os = socket.getOutputStream();
            InputStream is = socket.getInputStream();

            Scanner sc = new Scanner(System.in);
            
            System.out.println("\n<Menu>");
            System.out.println("Please input your client option in integer");
            System.out.println("Option 1) Convert text to upper case letters");
            System.out.println("Option 2) Ask the server what is my ip and port number");
            System.out.println("Option 3) Ask the server how many requests it has served so far");
            System.out.println("Option 4) Ask the server how long it has been running since it started");
            System.out.println("Option 5) Exit");
            System.out.print("Option :: ");
            
            int userOption = sc.nextInt();
            sc.nextLine();
            String messageCmd = null;
            int readByteCount;
            switch (userOption) {
                case 1:
                    messageCmd = "ASK_TXTCONV";
                    System.out.print("Input lowercase sentence: ");
                    String messageText = sc.nextLine();
                    os.write((messageCmd + "," + messageText).getBytes());
                    os.flush();
                    inputBuffer = new byte[1024];
                    readByteCount = is.read(inputBuffer);
                    recvMessage = new String(inputBuffer, 0, readByteCount, StandardCharsets.UTF_8);
                    System.out.println("Reply from server : "+ recvMessage);
                    is.close();
                    os.close();
                    socket.close();
                    sc.close();
                    break;
                case 2:
                    messageCmd = "ASK_IP_PORT";
                    os.write(messageCmd.getBytes());
                    os.flush();
                    inputBuffer = new byte[1024];
                    readByteCount = is.read(inputBuffer);
                    recvMessage = new String(inputBuffer, 0, readByteCount, StandardCharsets.UTF_8);
                    System.out.println("Reply from server : "+ recvMessage);
                    is.close();
                    os.close();
                    socket.close();
                    sc.close();
                    break;
                case 3:
                    messageCmd = "ASK_REQ_NUM";
                    os.write(messageCmd.getBytes());
                    os.flush();
                    inputBuffer = new byte[1024];
                    readByteCount = is.read(inputBuffer);
                    recvMessage = new String(inputBuffer, 0, readByteCount, StandardCharsets.UTF_8);
                    System.out.println("Reply from server : "+ recvMessage);
                    is.close();
                    os.close();
                    socket.close();
                    sc.close();
                    break;
                case 4:
                    messageCmd = "ASK_RUNTIME";
                    os.write(messageCmd.getBytes());
                    os.flush();
                    inputBuffer = new byte[1024];
                    readByteCount = is.read(inputBuffer);
                    recvMessage = new String(inputBuffer, 0, readByteCount, StandardCharsets.UTF_8);
                    System.out.println("Reply from server : requests served = "+ recvMessage);
                    is.close();
                    os.close();
                    socket.close();
                    sc.close();
                    break;
                case 5:
                    messageCmd = "ASK_CONNEND";
                    os.write(messageCmd.getBytes());
                    os.flush();

                    is.close();
                    os.close();
                    socket.close();
                    sc.close();
                    break;
                default:
                    break;
            }
            
        }
        catch (IOException e) {
            e.printStackTrace();
        }
        catch (NoSuchElementException e) {
            try {
                socket.close();
            }
            catch (IOException e1) {
                e1.printStackTrace();
            }
        }
        if (!socket.isClosed()) {
            try {
                socket.close();
            }
            catch (IOException e) {
                e.printStackTrace();
            }
        }
    }
}